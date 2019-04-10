import { BrowserModule }                                                                   from '@angular/platform-browser';
import { CheckResultComponent }                                                            from './check-result/check-result.component';
import { createCustomElement }                                                             from '@angular/elements';
import { ApplicationRef, DoBootstrap, Injector, NgModule }         from '@angular/core';
import { MatCardModule, MatTabsModule, MatExpansionModule, MatIconModule, MatChipsModule } from '@angular/material';
import { BrowserAnimationsModule }                                                         from '@angular/platform-browser/animations';
import { CommonModule }                                                                    from '@angular/common';

import { environment as env } from "../../../../src/environments/environment"

@NgModule({
  imports: [
    CommonModule,
    BrowserModule,
    MatCardModule,
    MatTabsModule,
    MatExpansionModule,
    MatIconModule,
    MatChipsModule,
    BrowserAnimationsModule,
  ],
  declarations: [
    CheckResultComponent,
  ],
  providers: [],
  entryComponents: [
    CheckResultComponent,
  ],
})
export class AppModule implements DoBootstrap {
  constructor(private injector: Injector) {
  }

  ngDoBootstrap(appRef: ApplicationRef): void {
    const nspCheckViewCE = createCustomElement(CheckResultComponent, { injector: this.injector });
    const webComponentName = "check-result"
    const selector = env.pluginName + "-" + webComponentName + "-" + env.pluginVersion;
    customElements.define(selector, nspCheckViewCE);

    const content = document.querySelector('head'); //TODO: move bus to document fragment
    content.dispatchEvent(new CustomEvent('CELoadedEvent', {
      detail: {
        pluginName: env.pluginName,
        pluginVersion: env.pluginVersion,
        webComponentName: webComponentName,
        selector: selector,
      },
      bubbles: false,
    }));
  }
}
