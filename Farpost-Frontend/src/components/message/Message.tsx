import MessageIcon from "../../assets/icons/attention.svg?react";
import CloseIcon from "../../assets/icons/x.svg?react";

interface Props {
  isClick: boolean;
  setIsClick: (isClick: boolean) => void;
}

function Message({ isClick, setIsClick }: Props) {

  return (
    isClick && (
      <div className="fixed top-0 left-0 w-full z-3 flex bg-message-bg min-h-12.5 items-center justify-center px-5">
        <div className="flex text-message-text font-bold underline gap-2.5 sm:gap-4.5 sm:justify-center justify-start items-center">
          <MessageIcon />
          <a className="w-[80%] min-w-[270px] sm:w-auto text-[14px]" href="#">
            Отопительный сезон во Владивостоке начинается 17 октября
          </a>
        </div>
        <button className="cursor-pointer" onClick={() => setIsClick(false)}>
          <CloseIcon className="absolute size-[30px] text-[#DEC47A] right-5 top-2.5 xl:right-25" />
        </button>
      </div>
    )
  );
}


export default Message;