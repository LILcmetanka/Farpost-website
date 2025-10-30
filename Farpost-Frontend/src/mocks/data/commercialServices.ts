import type { Service } from "../../types/service.types";

// Список из объектов сервисов (вода, отопление и т. д.)
export const commercialServices_1:Service[] = [
  { type: 'hot_water', count_buildings: 420, fraction_buildings: 0.86, time_last_blackout: '2019-10-30 11:15:00' },
  { type: 'cold_water', count_buildings: 37, fraction_buildings: 0.09, time_last_blackout: '2019-10-30 10:20:00' },
  { type: 'electricity', count_buildings: 115, fraction_buildings: 0.24, time_last_blackout: '2019-10-30 09:55:00' },
  { type: 'heat', count_buildings: 5, fraction_buildings: 0.01, time_last_blackout: '2019-10-29 10:25:00' },
];

export const commercialServices_2: Service[] = [
  { type: 'hot_water', count_buildings: 395, fraction_buildings: 0.81, time_last_blackout: '2019-10-31 09:10:00' },
  { type: 'cold_water', count_buildings: 42, fraction_buildings: 0.10, time_last_blackout: '2019-10-31 08:45:00' },
  { type: 'electricity', count_buildings: 97, fraction_buildings: 0.21, time_last_blackout: '2019-10-31 07:40:00' },
  { type: 'heat', count_buildings: 8, fraction_buildings: 0.02, time_last_blackout: '2019-10-30 10:05:00' },
];
