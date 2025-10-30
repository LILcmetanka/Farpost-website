export interface InfoOffs {
    service: string;
    start_off: string;
    end_off: string;
    amount_addresses: number;
}

export interface CalendarDay {
    date: string;
    services: string[];
}