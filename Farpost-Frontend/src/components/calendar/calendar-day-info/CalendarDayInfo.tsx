import type { CalendarDay, InfoOffs } from "../../../types/calendar.type";
import CircleIcon from "../../../assets/icons/circle.svg?react"
import { COLORMAP } from "../../../other-data/colorMap.data";
import { formatDateService } from "../../../utils/formatDateService";
import ChevronLeftIcon from "../../../assets/icons/chevron-left.svg?react"
import ChevronRightIcon from "../../../assets/icons/chevron-right.svg?react"
import { useEffect, useState } from "react";
import { NAMESSERVICES } from "../../../other-data/namesServices.data";

interface Props {
	day: CalendarDay | null;
}

function CalendarDayInfo({ day }: Props) {
	if (!day) return null;

	const [currentIndex, setCurrentIndex] = useState(0);
	const [blackouts, setBlackouts] = useState<InfoOffs[]>([]);

	useEffect(() => {
		async function fetchCalendarDay(selectDate: string): Promise<InfoOffs[]> {
			const baseUrl = import.meta.env.VITE_API_URL;
			let url = `${baseUrl}/off/calendar/day`;
			if (selectDate) {
				url += `?date=${encodeURIComponent(selectDate)}`;
			}

			const response = await fetch(url);
			if (!response.ok) {
				throw new Error(`Ошибка запроса ${response.status}`);
			}

			const data = await response.json();

      return Array.isArray(data) ? data : data.calendar?.calendar_day_1 ?? [];
		}

		fetchCalendarDay(day.date)
			.then((data) => {
				console.log("Отключения за день:", data);
				setBlackouts(data);
				setCurrentIndex(0);
			})
			.catch((err) => console.error("Ошибка загрузки:", err));
	}, [day.date]);

	const blackout = blackouts[currentIndex];

	const handlePrev = () => {
		setCurrentIndex((prev) => (prev === 0 ? blackouts.length - 1 : prev - 1));
	};

	const handleNext = () => {
		setCurrentIndex((prev) => (prev === blackouts.length - 1 ? 0 : prev + 1));
	};

	return (
		<div className="sm:relative w-[372px] sm:h-[150px] mb-10">
			<div className="sm:absolute inset-0 flex flex-col items-center justify-start gap-2.5">
				{blackouts.length === 0 || !blackout ? (
					<div className="text-center font-bold">Нет данных об отключениях</div>
				) : (
					<>
						<div className="flex flex-col gap-2.5 p-[15px] border border-line rounded-[10px] w-full">
							<h3 className="font-bold text-[16px]">{formatDateService(blackout.start_off)} – {formatDateService(blackout.end_off)}</h3>
							<div className="flex justify-between">
								<p>вид отключаемой услуги</p>
								<div className="flex items-center gap-2 w-[130px] justify-end">
									<CircleIcon className="size-[15px]" style={{ fill: COLORMAP[NAMESSERVICES[blackout.service]] }} />
									<p className="font-bold">{NAMESSERVICES[blackout.service]}</p>
								</div>
							</div>
							<div className="flex justify-between items-center">
								<p className="w-[180px]">количество адресов с отключением</p>
								<p className="font-bold">{blackout.amount_addresses}</p>
							</div>
							<a className="block self-end text-btn underline hover:text-btn-hover cursor-pointer" href="#">Подробнее</a>
						</div>

						{blackouts.length > 1 && (
							<div className="flex items-center gap-4 mt-2">
								<button onClick={handlePrev} className="hover:bg-line hover:rounded-full p-1 cursor-pointer">
									<ChevronLeftIcon />
								</button>
								<p className="font-bold">{currentIndex + 1} / {blackouts.length}</p>
								<button onClick={handleNext} className="hover:bg-line hover:rounded-full p-1 cursor-pointer">
									<ChevronRightIcon />
								</button>
							</div>
						)}
					</>
				)}
			</div>
		</div>
	);
}

export default CalendarDayInfo;
