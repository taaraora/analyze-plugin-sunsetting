import { Component, EventEmitter, Input, Output, ViewEncapsulation } from '@angular/core';
import { CheckResult, NodeCheckResult, WebComponentInfo } from "../models/models";

@Component({
  // it going to be generated while web component registering step
  // selector: 'app-check-result',
  templateUrl: './check-result.component.html',
  styleUrls: ['./check-result.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class CheckResultComponent {

  @Output('actionSubmit') actionSubmit = new EventEmitter<WebComponentInfo>();

  // event type === check-result
  @Input('checkResult')
  set checkResult(checkResult: string){
    this.currentCheck = JSON.parse(checkResult);
    let description: Array<NodeCheckResult> = JSON.parse(this.currentCheck.description as string);
    this.currentCheck.description = description;
  }

  public currentCheck: CheckResult;

  constructor() { }
}
