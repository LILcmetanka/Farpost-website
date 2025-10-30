import { useEffect, useState } from "react";
import LogIn from "../../assets/icons/log-in.svg?react";
import SearchField from "../elements/search-field/SearchField";

interface Props {
  isMessageVisible: boolean;
}

function Header({ isMessageVisible }: Props) {
  const [query, setQuery] = useState("");
  const [filteredStreets, setFilteredStreets] = useState<string[]>([]);
  const [isOpen, setIsOpen] = useState(false);

  useEffect(() => {
    async function fetchAndFilter() {
      try {
        const response = await fetch("/off/search");
        const data: string[] = await response.json();

        if (query.trim().length > 0) {
          const filtered = data.filter(street =>
            street.toLowerCase().includes(query.toLowerCase())
          );
          setFilteredStreets(filtered);
          setIsOpen(filtered.length > 0);
        } else {
          setFilteredStreets([]);
          setIsOpen(false);
        }
      } catch (err) {
        console.error(err);
        setFilteredStreets([]);
        setIsOpen(false);
      }
    }

    fetchAndFilter();
  }, [query]);

  return (
    <div className={`${isMessageVisible ? "mt-[70px]" : "mt-5"} flex justify-between items-center gap-10 mb-15 avg:mb-20 mx-5 xl:mx-25`}>
      <div className="flex-1 flex-col gap-2 justify-between w-[220px] flex avg:flex-row avg:items-center avg:justify-between avg:w-[60%]">
        <div className="flex items-end">
          <a className="block shrink-0" href="https://www.vl.ru/"><img className="w-[115px] h-auto" src="logo.svg" /></a>
          <a className="flex flex-col text-[18px] text-primary hover:text-primary-hover" href="#">
            <div>
              Отключение
              <span className="hidden sm:inline"> воды и света</span>
            </div>
            <span className="text-[12px] text-black">Владивостока</span>
          </a>
        </div>
        <div className="relative w-full max-w-[320px]">
          <SearchField value={query} onChange={setQuery} />

          {isOpen && filteredStreets.length > 0 && (
            <ul className="absolute z-2 w-full bg-white border border-line rounded-[10px] mt-1 shadow-lg max-h-[200px]">
              {filteredStreets.slice(0, 3).map((street) => (
                <li
                  key={street}
                  className="px-3 py-2 hover:bg-gray-100 cursor-pointer hover:rounded-[10px]"
                  onClick={() => {
                    setQuery(street);
                    setIsOpen(false);
                  }}
                >
                  {street}
                </li>
              ))}
            </ul>
          )}
        </div>

      </div>
      <div>
        <a className="flex gap-[5px] text-btn hover:text-btn-hover" href="#">
          войти
          <LogIn />
        </a>
      </div>
    </div>
  );
}

export default Header;