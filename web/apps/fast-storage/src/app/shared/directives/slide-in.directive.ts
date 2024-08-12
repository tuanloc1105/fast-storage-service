import { animate, AnimationBuilder, style } from '@angular/animations';
import { Directive, ElementRef, Input, OnInit, Renderer2 } from '@angular/core';

@Directive({
  selector: '[appSlideIn]',
  standalone: true,
})
export class SlideInDirective implements OnInit {
  @Input() appSlideIn: 'top' | 'left' | 'right' | 'bottom' = 'top';
  @Input() duration = '300ms';

  constructor(
    private el: ElementRef,
    private renderer: Renderer2,
    private builder: AnimationBuilder
  ) {}

  ngOnInit() {
    let initialStyles;
    let finalStyles;

    switch (this.appSlideIn) {
      case 'top':
        initialStyles = { transform: 'translateY(-100%)' };
        finalStyles = { transform: 'translateY(0%)' };
        break;
      case 'left':
        initialStyles = { transform: 'translateX(-100%)' };
        finalStyles = { transform: 'translateX(0%)' };
        break;
      case 'right':
        initialStyles = { transform: 'translateX(100%)' };
        finalStyles = { transform: 'translateX(0%)' };
        break;
      case 'bottom':
        initialStyles = { transform: 'translateY(100%)' };
        finalStyles = { transform: 'translateY(0%)' };
        break;
    }

    this.renderer.setStyle(
      this.el.nativeElement,
      'transform',
      initialStyles.transform
    );

    const animation = this.builder.build([
      style(initialStyles),
      animate(this.duration, style(finalStyles)),
    ]);

    const player = animation.create(this.el.nativeElement);
    player.play();
  }
}
