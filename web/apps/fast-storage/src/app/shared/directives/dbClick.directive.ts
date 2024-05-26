import {
  Directive,
  EventEmitter,
  HostListener,
  OnDestroy,
  OnInit,
  Output,
} from '@angular/core';
import { Subject, buffer, debounceTime, filter, map } from 'rxjs';

@Directive({
  selector: '[appDbClick]',
  standalone: true,
})
export class DbClickDirective implements OnInit, OnDestroy {
  private click$ = new Subject<MouseEvent>();

  @Output()
  doubleClick = new EventEmitter<MouseEvent>();

  @HostListener('click', ['$event'])
  onClick(event: MouseEvent) {
    this.click$.next(event);
  }

  ngOnInit() {
    this.click$
      .pipe(
        buffer(this.click$.pipe(debounceTime(250))),
        filter((list) => list.length === 2),
        map((list) => list[1])
      )
      .subscribe(this.doubleClick);
  }

  ngOnDestroy() {
    this.click$.complete();
  }
}
