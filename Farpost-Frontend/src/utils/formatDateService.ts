export function formatDateService(dateString: string): string {
    const date = new Date(dateString);
  
    const day = date.toLocaleString("ru-RU", { day: "2-digit" });
    const month = date.toLocaleString("ru-RU", { month: "short" }).replace(".", "");
    const hours = date.toLocaleString("ru-RU", { hour: "2-digit", minute: "2-digit" });
  
    return `${day} ${month}. ${hours}`;
  }
  