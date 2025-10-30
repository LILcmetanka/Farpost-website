import { useRef, useState, useEffect } from "react";
import ChevronDownIcon from "../../../assets/icons/chevron-down.svg?react";
import ChevronUpIcon from "../../../assets/icons/chevron-up.svg?react";
import type { Organization } from "../../../types/organization.types";
import type { Service } from "../../../types/service.types";
import { ICONMAP } from "../../../other-data/iconMap.data";
import { NAMESSERVICES } from "../../../other-data/namesServices.data";
import { COLORMAP } from "../../../other-data/colorMap.data";

interface Props {
  item: Service | Organization;
  activeItem: string | null;
  setActiveId: (value: null | string) => void;
}

function ServiceOrgsItem({ item, activeItem, setActiveId }: Props) {
  const { count_buildings, time_last_blackout } = item;

  let name: string;
  let specificField: string | number;
  let iconKey: string;

  if ("name" in item) {
    name = item.name;
    specificField = item.last_address;
    iconKey = "building";
  } else {
    name = NAMESSERVICES[item.type];
    specificField = item.fraction_buildings;
    iconKey = item.type;
  }

  const IconComponent = ICONMAP[iconKey];
  const isOpen = activeItem === name;

  const contentRef = useRef<HTMLDivElement>(null);
  const [height, setHeight] = useState(0);

  useEffect(() => {
    if (isOpen) {
      setHeight(contentRef.current?.scrollHeight || 0);
    } else {
      setHeight(0);
    }
  }, [isOpen]);

  return (
    <li className={`w-[332px] px-2.5 py-1.5 ${isOpen ? "bg-white rounded-[10px]" : "bg-auto"}`}>
      <button
        className="flex w-full text-left justify-between items-center cursor-pointer"
        onClick={() => setActiveId(isOpen ? null : name)}
      >
        <div className="flex items-center gap-5">
          <div><IconComponent /></div>
          <p className="uppercase leading-4">{name}</p>
        </div>
        <div className="size-6">
          {isOpen ? <ChevronUpIcon /> : <ChevronDownIcon />}
        </div>
      </button>

      <div ref={contentRef} className="transition-all duration-500 ease-in-out overflow-hidden" style={{ maxHeight: `${height}px`, opacity: isOpen ? 1 : 0 }}>
        <div className="flex flex-col gap-[15px]">
          <ul className="flex flex-col gap-2.5 px-[7px] py-2.5">
            <li className="flex justify-between">
              <div className="flex items-center gap-2.5">
                <span className="block size-[15px] bg-black rounded-full" style={{
    backgroundColor: COLORMAP[name] || "#0052CC",
  }}></span>
                <p>отсутствует</p>
              </div>
              <p className="font-bold w-[85px] text-right">{count_buildings} домов</p>
            </li>
            <li className="flex justify-between">
              <div className="flex items-center gap-2.5">
                <span className="block size-[15px] bg-black rounded-full" style={{
    backgroundColor: COLORMAP[name] || "#0052CC",
  }}></span>
                <p>{iconKey == "building" ? "адрес последнего отключения" : "охват от всех домов"}</p>
              </div>
              <p className="font-bold w-[85px] text-right whitespace-normal wrap-break-word">{specificField} { iconKey == "building" ? "" : "%"}</p>
            </li>
            <li className="flex justify-between">
              <div className="flex items-center gap-2.5">
                <span className="block size-[15px] bg-black rounded-full" style={{
    backgroundColor: COLORMAP[name] || "#0052CC",
  }}></span>
                <p>время последнего отключения</p>
              </div>
              <p className="font-bold w-[100px] text-right">{time_last_blackout}</p>
            </li>
          </ul>
          <a className="inline-block self-end text-btn hover:text-btn-hover underline" href="#">Подробнее</a>
        </div>
      </div>
    </li>
  );
}

export default ServiceOrgsItem;
