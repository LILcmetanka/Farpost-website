import CalendarBlock from "./calendar-block/CalendarBlock";
import CalendarLegends from "./calendar-legends/CalendarLegends";

function Calendar() {
  return (
    <div className="flex flex-col w-[372px] gap-[25px] mx-5 sm:ml-5 sm:mx-0 xl:ml-25">
      <h2 className="text-[21px] font-bold">Календарь плановых отключений</h2>
      <div className="flex flex-col gap-2.5">
        <CalendarLegends />
        <CalendarBlock />
      </div>
    </div>
  );
}

export default Calendar;