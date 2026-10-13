export const useScrollBottomListener = (bottomGap: number = 200) => {
  const isBottom = ref(false);

  const handleScroll = () => {
    const scrollPos = window.scrollY + window.innerHeight;
    const pageHeight = document.documentElement.scrollHeight;

    isBottom.value = scrollPos >= pageHeight - bottomGap;
  };

  onMounted(() => {
    window.addEventListener('scroll', handleScroll, { passive: true });
    handleScroll();
  });

  onBeforeUnmount(() => window.removeEventListener('scroll', handleScroll));

  return { isBottom };
};
