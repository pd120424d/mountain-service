import { HttpInterceptorFn, HttpRequest, HttpResponse, HttpHeaders } from '@angular/common/http';
import { of } from 'rxjs';
import { freshWindowInterceptor, FRESH_WINDOW_HEADER } from './fresh-window.interceptor';
import { TestBed } from '@angular/core/testing';

describe('freshWindowInterceptor', () => {
  beforeEach(() => {
    // No explicit reset of internal state; tests are designed to be independent within one spec
  });

  it('captures X-Fresh-Until from write response and attaches it to subsequent GETs', () => {
    const future = new Date(Date.now() + 60_000).toISOString();

    // Simulate a write that returns the header
    const writeReq = new HttpRequest('POST', '/api/v1/activities', {});
    const writeResp = new HttpResponse({ status: 200, headers: new HttpHeaders({ [FRESH_WINDOW_HEADER]: future }), body: {} });
    const writeNext = jasmine.createSpy('writeNext').and.returnValue(of(writeResp));

    TestBed.runInInjectionContext(() => {
      freshWindowInterceptor(writeReq, writeNext).subscribe();
    });

    // Next GET should receive the header attached
    const getReq = new HttpRequest('GET', '/api/v1/activities?pageSize=10');
    const getNext = jasmine.createSpy('getNext').and.returnValue(of(new HttpResponse({ status: 200 })));

    TestBed.runInInjectionContext(() => {
      freshWindowInterceptor(getReq, getNext).subscribe();
    });

    expect(getNext).toHaveBeenCalled();

    // Inspect the request passed to next: header should be present with the same value
    const calledWith = getNext.calls.mostRecent().args[0] as HttpRequest<any>;
    expect(calledWith.headers.get(FRESH_WINDOW_HEADER)).toBe(future);
  });

  it('does not attach the header to non-API GETs', () => {
    const getReq = new HttpRequest('GET', '/assets/logo.png');
    const getNext = jasmine.createSpy('getNext').and.returnValue(of(new HttpResponse({ status: 200 })));

    TestBed.runInInjectionContext(() => {
      freshWindowInterceptor(getReq, getNext).subscribe();
    });

    const calledWith = getNext.calls.mostRecent().args[0] as HttpRequest<any>;
    expect(calledWith.headers.has(FRESH_WINDOW_HEADER)).toBeFalse();
  });
});

