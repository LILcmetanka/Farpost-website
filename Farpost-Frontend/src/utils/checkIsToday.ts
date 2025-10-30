export function checkIsToday(datatString: string): boolean {
    const date = new Date(datatString);
    const now = new Date();

    return date.getDate() === now.getDate() && date.getMonth() === now.getMonth();
}