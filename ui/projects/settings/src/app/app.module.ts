import { Injector, NgModule, Provider, ValueProvider, ApplicationRef } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { createCustomElement } from "@angular/elements";

import { PluginSettingsComponent } from './plugin-settings/plugin-settings.component';
import { environment as env } from "../../../../src/environments/environment"

@NgModule({
  imports: [
    BrowserModule
  ],
  declarations: [
    PluginSettingsComponent
  ],
  providers: [],
  entryComponents: [
    PluginSettingsComponent
  ]
})
export class AppModule {
  constructor( private injector: Injector ) {  }

  ngDoBootstrap(appRef: ApplicationRef): void {
    const nspSettingsViewCE = createCustomElement(PluginSettingsComponent, { injector: this.injector });
    const webComponentName = "settings"
    const selector = env.pluginName + "-" + webComponentName + "-" + env.pluginVersion;
    customElements.define(selector, nspSettingsViewCE);

    const head = document.querySelector('head');
    head.dispatchEvent(new CustomEvent('CELoadedEvent', {
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
