import { HttpInterceptorFn, HttpRequest, HttpHandlerFn, HttpResponse } from '@angular/common/http';
import { tap } from 'rxjs/operators';

export const FRESH_WINDOW_HEADER = 'X-Fresh-Until';

let currentFresh = '';
let freshUntil = 0; // epoch ms

function updateFresh(v: string | null) {
  if (!v) return;
  if (!currentFresh || v > currentFresh) {
    currentFresh = v;
    const ts = Date.parse(currentFresh);
    freshUntil = Number.isNaN(ts) ? Date.now() + 2000 : ts; // fallback 2s if parse fails
  }
}

export const freshWindowInterceptor: HttpInterceptorFn = (req: HttpRequest<unknown>, next: HttpHandlerFn) => {
  const isWrite = req.method === 'POST' || req.method === 'PUT' || req.method === 'PATCH' || req.method === 'DELETE';
  const isRead = req.method === 'GET' && req.url.includes('/api/v1/');

  let outgoing = req;
  if (isRead && currentFresh && Date.now() < freshUntil) {
    outgoing = req.clone({ setHeaders: { [FRESH_WINDOW_HEADER]: currentFresh } });
  }

  return next(outgoing).pipe(
    tap(evt => {
      if (isWrite && evt instanceof HttpResponse) {
        const v = evt.headers.get(FRESH_WINDOW_HEADER);
        updateFresh(v);
      }
    })
  );
};

