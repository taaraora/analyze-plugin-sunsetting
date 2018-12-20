import {BrowserModule} from '@angular/platform-browser';
import {CheckResultComponent} from './check-result/check-result.component';
import {createCustomElement} from "@angular/elements";
import {settingsProvider} from "../../../settings/src/app/app.module";
import {Injector, NgModule} from "@angular/core";

export class ComponentSettings {
  constructor(public namePrefix: string) {}
}


@NgModule({
  imports: [
    BrowserModule
  ],
  declarations: [
    CheckResultComponent
  ],
  providers: [
    settingsProvider,
  ],
  entryComponents: [
    CheckResultComponent
  ]
})
export class AppModule {
  constructor(private injector: Injector) {
  }

  ngDoBootstrap() {
    const nspSettingsVewCE = createCustomElement(CheckResultComponent, {injector: this.injector});
    //TODO: add web component selector generation mechanism
    customElements.define("check-result", nspSettingsVewCE);

    const content = document.querySelector('head');
    content.dispatchEvent(new CustomEvent('loadingNotifier', {
      detail: {
        pluginName: 'analyze-plugin-sunsetting',
        pluginVersion: 'v2.0.0',
        webComponentName: 'check-result',
        selector: 'analyze-plugin-sunsetting-check-result-v2-0-0',
      },
      bubbles: true
    }));
    console.log('loadingNotifier emitted main')
  }
}
