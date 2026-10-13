const colorClasses = [
  'bg-red-500',
  'bg-orange-500',
  'bg-lime-500',
  'bg-green-500',
  'bg-emerald-500',
  'bg-sky-500',
  'bg-blue-500',
  'bg-purple-500',
  'bg-yellow-500',
  'bg-red-500',
];

export function textLogo(text?: string): { text: string; class: string } {
  const sliceText = text ? text.slice(0, 3) : 'none';
  const charCode = sliceText.charAt(0).toLowerCase().charCodeAt(0);
  const colorIndex = charCode % colorClasses.length;
  return {
    text: sliceText,
    class: colorClasses[colorIndex] || 'bg-gray-500',
  };
}

export function lightenColorWithOpacity(hex: string | null | undefined, lightenPercent = 20, opacity = 0.7) {
  if (!hex) return `rgba(0,0,0,0.1)`;

  hex = hex.replace(/^#/, '');

  const r = parseInt(hex.substring(0, 2), 16);
  const g = parseInt(hex.substring(2, 4), 16);
  const b = parseInt(hex.substring(4, 6), 16);

  const lightenFactor = lightenPercent / 100;
  const newR = Math.round(r + (255 - r) * lightenFactor);
  const newG = Math.round(g + (255 - g) * lightenFactor);
  const newB = Math.round(b + (255 - b) * lightenFactor);

  return `rgba(${newR}, ${newG}, ${newB}, ${opacity})`;
}
