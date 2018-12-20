import {AfterViewInit, Component, EventEmitter, Input, OnInit, Output, ViewEncapsulation} from '@angular/core';

@Component({
  // selector: 'app-check-result',
  templateUrl: './check-result.component.html',
  styleUrls: ['./check-result.component.css'],
  encapsulation: ViewEncapsulation.Emulated
})
export class CheckResultComponent implements AfterViewInit {

  @Output('actionSubmit') actionSubmit = new EventEmitter<WebComponentInfo>();

  @Input('checkResult')
  set checkResult(result: string) {
    console.debug('client-a received state at plugin', result);
  }

  constructor() { }

  ngAfterViewInit(): void {
    setTimeout( () => {
      this.actionSubmit.emit(new class implements WebComponentInfo {
        pluginName: 'analyze-plugin-sunsetting';
        pluginVersion: 'v2.0.0';
        webComponentName: 'check-result';
        selector: 'analyze-plugin-sunsetting-check-result-v2-0-0';
      });
      console.log('loadingNotifier emitted from plugin')
    }, 2000);
  }


}


export interface WebComponentInfo {
  selector: string;
  webComponentName: string;
  pluginName: string;
  pluginVersion: string;
}
