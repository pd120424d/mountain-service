import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LocalDatePipe, LocalDateTimePipe } from './date-time.pipe';

@NgModule({
  imports: [CommonModule, LocalDatePipe, LocalDateTimePipe],
  exports: [LocalDatePipe, LocalDateTimePipe]
})
export class DateTimeModule {}

