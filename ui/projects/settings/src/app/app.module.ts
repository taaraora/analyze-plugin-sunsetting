import { BrowserModule } from '@angular/platform-browser';
import {Injector, NgModule, Provider, ValueProvider} from '@angular/core';

import { PluginSettingsComponent } from './plugin-settings/plugin-settings.component';
import {createCustomElement} from "@angular/elements";

export let settingsProvider: ValueProvider = null;

export const InitSettingProvider = (settings: ComponentSettings) => {
  settingsProvider = {
    provide: ComponentSettings,
    useValue: settings,
  };
};

export class ComponentSettings {
  constructor(public namePrefix: string) {}
}


@NgModule({
  imports: [
    BrowserModule
  ],
  declarations: [
    PluginSettingsComponent
  ],
  providers: [
    settingsProvider,
  ],
  entryComponents: [
    PluginSettingsComponent
  ]
})
export class AppModule {
  constructor(private injector: Injector, private componentSettings: ComponentSettings) {
  }

  ngDoBootstrap() {
    const nspSettingsVewCE = createCustomElement(PluginSettingsComponent, { injector: this.injector });
    customElements.define(this.componentSettings.namePrefix, nspSettingsVewCE);
  }
}
