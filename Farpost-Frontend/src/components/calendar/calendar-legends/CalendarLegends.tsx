import { COLORMAP } from "../../../other-data/colorMap.data";

function CalendarLegends() {
	const namesLegends = [
		{id: 1, title: "горячая вода"},
		{id: 2, title: "электричество"},
		{id: 3, title: "отопление"},
		{id: 4, title: "холодная вода"},
	];

	return (
		<div>
			<ul className="flex text-[12px] justify-center gap-x-2 px-[35px] flex-wrap">
				{namesLegends.map((item) => (
					<li key={item.id} className="flex items-center gap-0.5">
						<span className="block size-2.5" style={{ backgroundColor: COLORMAP[item.title] }}></span>
						<p>{item.title}</p>
					</li>
				))}
			</ul>
		</div>
	);
}

export default CalendarLegends;