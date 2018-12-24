import {BrowserModule} from '@angular/platform-browser';
import {CheckResultComponent} from './check-result/check-result.component';
import {createCustomElement} from "@angular/elements";
import {ApplicationRef, CUSTOM_ELEMENTS_SCHEMA, DoBootstrap, Injector, NgModule} from "@angular/core";
import { MatCardModule, MatTabsModule, MatExpansionModule, MatIconModule, MatChipsModule } from "@angular/material";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { CommonModule } from '@angular/common';

@NgModule({
  imports: [
    CommonModule,
    BrowserModule,
    MatCardModule,
    MatTabsModule,
    MatExpansionModule,
    MatIconModule,
    MatChipsModule,
    BrowserAnimationsModule
  ],
  declarations: [
    CheckResultComponent
  ],
  providers: [],
  entryComponents: [
    CheckResultComponent
  ]
})
export class AppModule implements DoBootstrap {
  constructor(private injector: Injector) {
  }

  ngDoBootstrap(appRef: ApplicationRef): void {
    const nspSettingsVewCE = createCustomElement(CheckResultComponent, {injector: this.injector});
    //TODO: add web component selector generation mechanism
    customElements.define("analyze-plugin-sunsetting-check-result-v2-0-0", nspSettingsVewCE);

    const content = document.querySelector('head'); //TODO: move bus to document fragment
    content.dispatchEvent(new CustomEvent('loadingNotifier', {
      detail: {
        pluginName: 'analyze-plugin-sunsetting',
        pluginVersion: 'v2.0.0',
        webComponentName: 'check-result',
        selector: 'analyze-plugin-sunsetting-check-result-v2-0-0',
      },
      bubbles: true
    }));
    console.debug('loadingNotifier emitted from plugin app module')
  }
}
