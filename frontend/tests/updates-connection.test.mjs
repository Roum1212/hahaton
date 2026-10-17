import assert from 'node:assert/strict';
import { setImmediate as flush, setTimeout as sleep } from 'node:timers/promises';
import { test } from 'node:test';
import { Code, ConnectError } from '@connectrpc/connect';
import { loadTypeScript } from './helpers/load-typescript.mjs';

const UpdateKind = { READY: 1, CHANGED: 2, HEARTBEAT: 3 };

const setup = (t) => {
  const streams = [];
  let expired = 0;
  const store = {
    connected: false,
    revision: 0,
    notify() {
      this.revision += 1;
    },
    service: {
      subscribe(signal) {
        let resolve;
        let reject;
        const stream = {
          signal,
          send: (value) => resolve({ value, done: false }),
          fail: (error) => reject(error),
          [Symbol.asyncIterator]() {
            return this;
          },
          next() {
            return new Promise((yes, no) => {
              resolve = yes;
              reject = no;
            });
          },
        };
        streams.push(stream);
        return stream;
      },
    },
  };
  const { useUpdatesConnection } = loadTypeScript(
    new URL('../composables/updates/useUpdates.ts', import.meta.url),
    {
      AbortController,
      setTimeout,
      clearTimeout,
      useUpdatesStore: () => store,
      useAuthStore: () => ({
        expire: () => {
          expired += 1;
        },
      }),
    },
    { '~/core/updates/services/updates/updates.types': { UpdateKind } },
  );
  const connection = useUpdatesConnection();
  t.after(connection.stop);
  return { connection, store, streams, expired: () => expired };
};

test('READY refreshes data loaded before the first subscription', async (t) => {
  const { connection, store, streams } = setup(t);
  connection.start();
  streams[0].send(UpdateKind.READY);
  await sleep(350);
  assert.equal(store.revision, 1);
});

test('login can start a new stream after an unauthenticated response', async (t) => {
  const { connection, streams, expired } = setup(t);
  connection.start();
  streams[0].fail(new ConnectError('expired session', Code.Unauthenticated));
  await flush();
  assert.equal(expired(), 1);

  connection.start();
  assert.equal(streams.length, 2);
});

test('an old canceled stream cannot disconnect the new login', async (t) => {
  const { connection, store, streams } = setup(t);
  connection.start();
  connection.stop();
  connection.start();
  streams[1].send(UpdateKind.READY);
  await flush();
  assert.equal(store.connected, true);

  streams[0].fail(new ConnectError('request canceled', Code.Canceled));
  await flush();
  assert.equal(store.connected, true);
});

test('continuous changes do not postpone refresh indefinitely', async (t) => {
  const { connection, store, streams } = setup(t);
  connection.start();
  for (let i = 0; i < 5; i += 1) {
    streams[0].send(UpdateKind.CHANGED);
    await sleep(100);
  }
  assert.ok(store.revision > 0);
});

test('temporary network failure reconnects and synchronizes data', async (t) => {
  const { connection, store, streams } = setup(t);
  connection.start();
  streams[0].fail(new ConnectError('network unavailable', Code.Unavailable));
  await sleep(1100);
  assert.equal(streams.length, 2);

  streams[1].send(UpdateKind.READY);
  await sleep(350);
  assert.equal(store.connected, true);
  assert.equal(store.revision, 1);
});

test('heartbeat keeps the connection alive without refreshing data', async (t) => {
  const { connection, store, streams } = setup(t);
  connection.start();
  connection.start();
  assert.equal(streams.length, 1);

  streams[0].send(UpdateKind.HEARTBEAT);
  await sleep(350);
  assert.equal(store.connected, true);
  assert.equal(store.revision, 0);
});
