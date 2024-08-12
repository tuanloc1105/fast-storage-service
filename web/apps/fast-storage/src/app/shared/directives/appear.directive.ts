import { animate, AnimationBuilder, style } from '@angular/animations';
import { Directive, ElementRef, Input, OnInit, Renderer2 } from '@angular/core';

@Directive({
  selector: '[appAppear]',
  standalone: true,
})
export class AppearDirective implements OnInit {
  @Input() duration = '200ms';
  @Input() delay = '0ms';
  @Input() moveDistance = '10px';

  constructor(
    private el: ElementRef,
    private renderer: Renderer2,
    private builder: AnimationBuilder
  ) {}

  ngOnInit() {
    // Initially set the opacity to 0
    this.renderer.setStyle(this.el.nativeElement, 'opacity', '0');
    this.renderer.setStyle(
      this.el.nativeElement,
      'transform',
      `translateY(${this.moveDistance})`
    );

    // Define the animation: fade in and move to its original position
    const animation = this.builder.build([
      style({ opacity: 0, transform: `translateY(${this.moveDistance})` }),
      animate(
        `${this.duration} ${this.delay}`,
        style({ opacity: 1, transform: 'translateY(0)' })
      ),
    ]);

    // Create the animation player and play it
    const player = animation.create(this.el.nativeElement);
    player.play();
  }
}
