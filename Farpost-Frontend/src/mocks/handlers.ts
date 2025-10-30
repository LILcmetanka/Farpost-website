import { http, HttpResponse } from 'msw';
import { data } from './data/data';

export const handlers = [
  http.get('/off/search', () => {
    return HttpResponse.json(data.streetsData.streets_1);
  }),
  
  http.get('/off/blackouts', () => {
    return HttpResponse.json(data.commercialServices.commercialServices_1);
  }),

  http.get('/off/orgs', () => {
    return HttpResponse.json(data.organizations.organizations_1);
  }),

  http.get('/off/complaints', (req) => {
    const url = new URL(req.request.url);
    const period = url.searchParams.get('period')?.trim().toLowerCase();

    console.log('Request URL:', req.request.url);
    console.log('period:', period);

    switch (period) {
      case 'hour':
        return HttpResponse.json(data.complaintsHour.complaintsHour_1);
      case 'day':
        return HttpResponse.json(data.complaintsDay.complaintsDay_1);
      case 'week':
        return HttpResponse.json(data.complaintsWeek.complaintsWeek_1);
      case 'month':
        return HttpResponse.json(data.complaintsMonth.complaintsMonth_1);
      default:
        return HttpResponse.json([]);
    }
  }),

  http.get('/off/calendar', (req) => {
    const url = new URL(req.request.url);
    const month = url.searchParams.get('month')?.trim().toLowerCase();

    if (month === '2019-10') {
      return HttpResponse.json(data.calendar.calendar_1);
    }
    
    return HttpResponse.json([]);
  }),

  http.get('/off/calendar/day', (req) => {
    const url = new URL(req.request.url);
    const date = url.searchParams.get('date')?.trim().toLowerCase();

    if (date === '2019-10-01') {
      return HttpResponse.json(data.calendar.calendar_day_1);
    }

    return HttpResponse.json([]);
  }),

];
