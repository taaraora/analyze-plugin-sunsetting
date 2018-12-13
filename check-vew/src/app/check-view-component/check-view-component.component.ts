import {Component, Input, OnInit, ViewEncapsulation} from '@angular/core';

@Component({
  // selector is not needed  because CustomElement gets one assigned when it is registered
  // selector: 'app-check-view-component',
  templateUrl: './check-view-component.component.html',
  styleUrls: ['./check-view-component.component.scss'],
  encapsulation: ViewEncapsulation.Emulated
})
export class CheckViewComponentComponent implements OnInit {

  @Input() a: number;
  @Input() b: number;
  @Input() c: number;

  constructor() { }

  ngOnInit() {
  }

}
