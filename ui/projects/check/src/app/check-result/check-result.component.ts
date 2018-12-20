import {Component, EventEmitter, OnInit, Output, ViewEncapsulation} from '@angular/core';

@Component({
  // selector: 'app-check-result',
  templateUrl: './check-result.component.html',
  styleUrls: ['./check-result.component.css'],
  encapsulation: ViewEncapsulation.ShadowDom
})
export class CheckResultComponent implements OnInit {

  constructor() { }

  ngOnInit() {
    this.loadingNotifier.emit(new class implements WebComponentInfo {
      pluginName: 'analyze-plugin-sunsetting';
      pluginVersion: 'v2.0.0';
      webComponentName: 'check-result';
      selector: 'analyze-plugin-sunsetting-check-result-v2-0-0';
    });
    console.log('loadingNotifier emitted')
  }

  @Output() loadingNotifier = new EventEmitter<WebComponentInfo>();

}


export interface WebComponentInfo {
  selector: string;
  webComponentName: string;
  pluginName: string;
  pluginVersion: string;
}
