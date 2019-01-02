import {AfterViewInit, Component, EventEmitter, Input, Output, ViewEncapsulation} from '@angular/core';
import {CheckResult, WebComponentInfo} from "../models/models";

@Component({
  // it going to be generated while web component registering step
  // selector: 'app-check-result',
  templateUrl: './check-result.component.html',
  styleUrls: ['./check-result.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class CheckResultComponent implements AfterViewInit {

  @Output('actionSubmit') actionSubmit = new EventEmitter<WebComponentInfo>();

  // event type === check-result
  @Input('checkResult')
  set checkResult(checkResult: string){
    this.currentCheck = JSON.parse(checkResult);
    console.info(this.currentCheck);
  }

  private currentCheck: CheckResult;

  constructor() { }

  ngAfterViewInit(): void {
    setTimeout( () => {
      this.actionSubmit.emit( {
        pluginName: 'analyze-plugin-sunsetting',
        pluginVersion: 'v2.0.0',
        webComponentName: 'check-result',
        selector: 'analyze-plugin-sunsetting-check-result-v2-0-0',
      });
      console.log('loadingNotifier emitted from plugin')
    }, 2000);
  }


}
