import {enableProdMode} from '@angular/core';
import {platformBrowserDynamic} from '@angular/platform-browser-dynamic';

import {AppModule, InitSettingProvider, ComponentSettings} from './app/app.module';
import {environment} from './environments/environment';


if (environment.production) {
  enableProdMode();
}


export function main(customElementTag: string) {
  InitSettingProvider(new ComponentSettings(customElementTag));
  platformBrowserDynamic().bootstrapModule(AppModule).catch(err => console.error(err));
}

