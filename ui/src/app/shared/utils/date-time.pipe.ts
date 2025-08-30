import { Pipe, PipeTransform } from '@angular/core';

@Pipe({ name: 'localDateTime', standalone: true })
export class LocalDateTimePipe implements PipeTransform {
  transform(value?: string | Date, format: Intl.DateTimeFormatOptions = {}): string {
    if (!value) return '';
    try {
      const date = value instanceof Date ? value : new Date(value);
      return new Intl.DateTimeFormat(undefined, {
        year: 'numeric', month: 'short', day: '2-digit',
        hour: '2-digit', minute: '2-digit', hour12: false,
        ...format,
      }).format(date);
    } catch {
      return String(value);
    }
  }
}

@Pipe({ name: 'localDate', standalone: true })
export class LocalDatePipe implements PipeTransform {
  transform(value?: string | Date, format: Intl.DateTimeFormatOptions = {}): string {
    if (!value) return '';
    try {
      const date = value instanceof Date ? value : new Date(value);
      return new Intl.DateTimeFormat(undefined, {
        year: 'numeric', month: 'short', day: '2-digit',
        ...format,
      }).format(date);
    } catch {
      return String(value);
    }
  }
}

