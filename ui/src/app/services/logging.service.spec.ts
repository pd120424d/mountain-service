import { TestBed } from "@angular/core/testing";
import { LoggingService } from "./logging.service";

describe('LoggingService', () => {
  let service: LoggingService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [LoggingService]
    });
    service = TestBed.inject(LoggingService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should log info message', () => {
    const spy = spyOn(console, 'info');
    service.info('Test info message');
    expect(spy).toHaveBeenCalledWith('[INFO] Test info message', '');
  });

  it('should log warn message', () => {
    const spy = spyOn(console, 'warn');
    service.warn('Test warn message');
    expect(spy).toHaveBeenCalledWith('[WARN] Test warn message', '');
  });

  it('should log error message', () => {
    const spy = spyOn(console, 'error');
    service.error('Test error message');
    expect(spy).toHaveBeenCalledWith('[ERROR] Test error message', '');
  });

  it('should log debug message', () => {
    const spy = spyOn(console, 'debug');
    service.debug('Test debug message');
    expect(spy).toHaveBeenCalledWith('[DEBUG] Test debug message', '');
  });

  it('should log info message with data', () => {
    const spy = spyOn(console, 'info');
    const testData = { key: 'value' };
    service.info('Test info message', testData);
    expect(spy).toHaveBeenCalledWith('[INFO] Test info message', testData);
  });

  it('should log warn message with data', () => {
    const spy = spyOn(console, 'warn');
    const testData = { error: 'warning' };
    service.warn('Test warn message', testData);
    expect(spy).toHaveBeenCalledWith('[WARN] Test warn message', testData);
  });

  it('should log error message with data', () => {
    const spy = spyOn(console, 'error');
    const testData = { stack: 'error stack' };
    service.error('Test error message', testData);
    expect(spy).toHaveBeenCalledWith('[ERROR] Test error message', testData);
  });

  it('should log debug message with data', () => {
    const spy = spyOn(console, 'debug');
    const testData = { debug: 'info' };
    service.debug('Test debug message', testData);
    expect(spy).toHaveBeenCalledWith('[DEBUG] Test debug message', testData);
  });
});