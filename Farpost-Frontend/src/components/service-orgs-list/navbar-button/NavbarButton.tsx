interface Props {
  name: string;
  isActive?: boolean;
  onClick: () => void;
}

function NavbarButton({ name, isActive = false, onClick }: Props) {
  return (
    <button
      className={`${isActive ? "bg-primary text-white" : "bg-black/1 border border-line"} font-bold p-2.5 rounded-[10px] cursor-pointer`}
      onClick={onClick}
    >
      {name}
    </button>
  );
}
export default NavbarButton;