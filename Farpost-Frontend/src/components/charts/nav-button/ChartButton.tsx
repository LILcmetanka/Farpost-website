interface Props {
  title: string;
  active?: boolean;
  onClick: () => void;
}

function ChartButton({ title, active = false, onClick }: Props) {
  return (
    <button className={`cursor-pointer relative pb-2  transition-colors duration-300 ${active && "text-primary"}`} onClick={onClick}>
      {title}
      {active &&
        <span className="absolute bottom-0 left-1/2 -translate-x-1/2 w-full h-[5px] bg-primary rounded-t-[10px]" />
      }
    </button>
  );
}

export default ChartButton;