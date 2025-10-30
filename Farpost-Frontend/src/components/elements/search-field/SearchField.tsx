import SearchIcon from "../../../assets/icons/search.svg?react";

interface Props {
  value: string;
  onChange: (value: string) => void;
}

function SearchField({ value, onChange }: Props) {
  return (
    <div className="group">
      <label className="flex items-center gap-2 px-3.5 py-2.5 w-full max-w-[320px] min-w-[220px] border border-gray-300 group-focus-within:border-primary rounded-lg duration-300 sm:w-[320px] group-focus-within:shadow-[0_0_10px_3px_#FFF0E8]">
        <SearchIcon className="w-5 h-5 text-gray-500 duration-300 group-focus-within:text-primary" />
        <input
          className="outline-none text-[16px] w-full bg-transparent [&::-webkit-search-cancel-button]:appearance-none placeholder:text-gray-500"
          type="search"
          value={value}
          placeholder="Введите адрес"
          onChange={(e) => onChange(e.target.value)}
        />
      </label>
    </div>
  );
}

export default SearchField;