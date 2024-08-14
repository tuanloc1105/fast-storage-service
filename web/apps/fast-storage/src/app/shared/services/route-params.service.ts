import { Injectable } from '@angular/core';
import { ActivatedRoute, Params, Router } from '@angular/router';

@Injectable({
  providedIn: 'root',
})
export class RouteParamsService {
  constructor(private router: Router, private activatedRoute: ActivatedRoute) {}

  public getParams(): Params {
    return this.activatedRoute.snapshot.params;
  }

  public getParam(param: string): string {
    return this.activatedRoute.snapshot.params[param];
  }

  public setRouteParams(params: Params): void {
    this.router.navigate([], {
      relativeTo: this.activatedRoute,
      queryParams: params,
      queryParamsHandling: 'merge',
    });
  }
}
