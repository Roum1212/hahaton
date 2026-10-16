import { LANDING_CITIES, type ILandingCity } from '~/config/landing.config';

export const useCity = () => {
  const city = useLocalStorage('city', LANDING_CITIES[0]!.name);
  const confirmed = useLocalStorage('cityConfirmed', false);
  const dialogVisible = useState('cityDialogVisible', () => false);

  const current = computed<ILandingCity>(() => LANDING_CITIES.find((item) => item.name === city.value) ?? LANDING_CITIES[0]!);

  const confirm = () => {
    confirmed.value = true;
  };

  const openDialog = () => {
    confirmed.value = true;
    dialogVisible.value = true;
  };

  const select = (value: string) => {
    city.value = value;
    confirmed.value = true;
    dialogVisible.value = false;
  };

  return { city, current, cities: LANDING_CITIES, confirmed, dialogVisible, confirm, openDialog, select };
};
