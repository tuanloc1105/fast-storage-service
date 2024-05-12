import {
  HttpErrorResponse,
  HttpInterceptorFn,
  HttpResponse,
} from '@angular/common/http';
import { inject } from '@angular/core';

// #docregion excerpt
import { ElectronService, LoggingService } from '@app/shared/services';
import { finalize, tap } from 'rxjs';

export const loggingInterceptor: HttpInterceptorFn = (req, next) => {
  const logging = inject(LoggingService);
  const electron = inject(ElectronService);

  const started = Date.now();
  let ok: string;
  let description: string;

  return next(req).pipe(
    tap({
      next: (event) => {
        if (event instanceof HttpResponse) {
          ok = 'succeeded';
          description = (event.body as any).errorMessage;
        }
      },
      error: (error: HttpErrorResponse) => {
        (ok = 'failed'), (description = error.error.errorMessage);
      },
    }),
    finalize(() => {
      const elapsed = Date.now() - started;
      const msg = `${req.method} "${req.urlWithParams}"
          Status: ${ok} in ${elapsed} ms.
          Description: ${description}`;
      logging.add(msg);
      electron.log(msg);
    })
  );
};
