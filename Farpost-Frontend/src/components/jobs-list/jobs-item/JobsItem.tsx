interface Props {
  jobName: string;
  salary: string
}

function JobsItem({ jobName, salary }: Props) {
  return (
    <li className="group">
      <a className="flex justify-between items-center w-[342px] px-[15px] py-2.5 bg-[#F6FAFE] text-[12px] font-bold text-[#0F71C8] leading-4" href="#">
        <p className="w-[150px] group-hover:text-btn-hover">{jobName}</p>
        <div className="bg-[#E6F1FD] text-[#0A4E8B] p-2.5 rounded-[10px] whitespace-nowrap">
          <p className="tracking-normal">{salary}</p>
        </div>
      </a>
    </li>
  );
}

export default JobsItem;        