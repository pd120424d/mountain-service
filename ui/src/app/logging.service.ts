import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class LoggingService {
  info(message: string, data?: any): void {
    console.info(`[INFO] ${message}`, data || '');
  }

  warn(message: string, data?: any): void {
    console.warn(`[WARN] ${message}`, data || '');
  }

  error(message: string, data?: any): void {
    console.error(`[ERROR] ${message}`, data || '');
  }

  debug(message: string, data?: any): void {
    console.debug(`[DEBUG] ${message}`, data || '');
  }
}
