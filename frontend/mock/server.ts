import http from 'node:http';
import type { ConnectRouter, Interceptor } from '@connectrpc/connect';
import { connectNodeAdapter } from '@connectrpc/connect-node';
import { MAX_IMAGE_BYTES, serveStorage } from './images';
import { seed } from './seed';
import { registerAdminService } from './services/admin';
import { registerCooperationService } from './services/cooperation';
import { registerCoordinatorService } from './services/coordinator';
import { registerEventChatService } from './services/event-chat';
import { registerEventService } from './services/event';
import { registerGeoService } from './services/geo';
import { registerImageService } from './services/image';
import { registerParticipationService } from './services/participation';
import { registerRewardOrderService } from './services/reward-order';
import { registerRewardService } from './services/reward';
import { registerSessionService } from './services/session';
import { registerUpdatesService } from './services/updates';
import { registerVolunteerService } from './services/volunteer';
import { updatesHub } from './updates-hub';

const PORT = Number(process.env.MOCK_PORT || 8080);
const ADMIN_PORT = Number(process.env.MOCK_ADMIN_PORT || 8081);

const clientRoutes = (router: ConnectRouter) => {
  registerSessionService(router);
  registerVolunteerService(router);
  registerCoordinatorService(router);
  registerCooperationService(router);
  registerEventService(router);
  registerParticipationService(router);
  registerRewardService(router);
  registerRewardOrderService(router);
  registerEventChatService(router);
  registerGeoService(router);
  registerImageService(router);
  registerUpdatesService(router);
};

const READ_METHOD_PREFIXES = ['Get', 'List', 'Find', 'Suggest', 'Subscribe'];
const SILENT_SERVICES = ['session.v1.SessionService', 'image.v1.ImageService'];

const publishUpdates: Interceptor = (next) => async (req) => {
  const response = await next(req);
  const read = READ_METHOD_PREFIXES.some((prefix) => req.method.name.startsWith(prefix));
  if (!req.stream && !read && !SILENT_SERVICES.includes(req.service.typeName)) updatesHub.notifyAll();
  return response;
};

const READ_MAX_BYTES = 8 * 1024 * 1024;

const withCors =
  (handler: http.RequestListener): http.RequestListener =>
  (req, res) => {
    const origin = req.headers.origin;
    if (origin) {
      res.setHeader('Access-Control-Allow-Origin', origin);
      res.setHeader('Access-Control-Allow-Credentials', 'true');
      res.setHeader('Vary', 'Origin');
    }
    if (req.method === 'OPTIONS') {
      res.setHeader('Access-Control-Allow-Methods', 'POST, GET, OPTIONS');
      res.setHeader('Access-Control-Allow-Headers', req.headers['access-control-request-headers'] || '*');
      res.setHeader('Access-Control-Max-Age', '7200');
      res.writeHead(204).end();
      return;
    }
    if (serveStorage(req, res)) return;
    handler(req, res);
  };

const invitations = seed();

http
  .createServer(
    withCors(
      connectNodeAdapter({
        routes: clientRoutes,
        interceptors: [publishUpdates],
        readMaxBytes: Math.max(READ_MAX_BYTES, MAX_IMAGE_BYTES * 2),
      }),
    ),
  )
  .listen(PORT, () => {
    console.log(`Mock Connect RPC: http://localhost:${PORT}`);
    console.log(`DaData: ${process.env.MOCK_DADATA_API_KEY ? 'подключён' : 'не настроен (MOCK_DADATA_API_KEY)'}`);
    invitations.forEach((item) => console.log(`Приглашение от ${item.coordinator}: ${item.url}`));
  });

http.createServer(connectNodeAdapter({ routes: registerAdminService, interceptors: [publishUpdates] })).listen(ADMIN_PORT, () => {
  console.log(`Mock AdminService: http://localhost:${ADMIN_PORT} (Authorization: Bearer $MOCK_ADMIN_TOKEN)`);
});
