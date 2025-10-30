import { DATAJOBS } from "../../other-data/job.data";
import JobsItem from "./jobs-item/JobsItem";

function JobsList() {

  return (
    <div className="flex flex-col max-h-[290px] h-full gap-[25px] px-[15px] py-[25px] max-w-[372px] border border-line rounded-[10px] mx-5 sm:mr-5 sm:mx-0 xl:mr-25">
      <h3 className="px-[15px] text-[14px] font-bold underline text-btn hover:text-btn-hover"><a href="#">Работа с высоким доходом</a></h3>
      <ul className="flex flex-col gap-3">
        {DATAJOBS.map((item) => (
          <JobsItem key={item.id} jobName={item.name} salary={item.salary} />
        ))}
      </ul>
    </div>
  );
}

export default JobsList;