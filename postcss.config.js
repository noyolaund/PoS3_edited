/*
    Un archivo vacío para que no genere un error que dice:

    ERROR in ./node_modules/vuetify/dist/vuetify.min.css
    Module build failed: ModuleBuildError: Module build failed: Error: No PostCSS Config found in: /home/parzibyte/apps/node_modules/vuetify/dist
        at config.load.then (/home/parzibyte/apps/node_modules/postcss-load-config/src/index.js:55:15)
        at <anonymous>
        at runLoaders (/home/parzibyte/apps/node_modules/webpack/lib/NormalModule.js:195:19)
        at /home/parzibyte/apps/node_modules/loader-runner/lib/LoaderRunner.js:364:11
        at /home/parzibyte/apps/node_modules/loader-runner/lib/LoaderRunner.js:230:18
        at context.callback (/home/parzibyte/apps/node_modules/loader-runner/lib/LoaderRunner.js:111:13)
        at Promise.resolve.then.then.catch (/home/parzibyte/apps/node_modules/postcss-loader/lib/index.js:194:71)
        at <anonymous>
    @ ./src/main.js 7:0-38

    ERROR in ./node_modules/vuetify/dist/vuetify.min.css
    Module build failed: ModuleBuildError: Module build failed: Error: No PostCSS Config found in: /home/parzibyte/apps/node_modules/vuetify/dist
        at config.load.then (/home/parzibyte/apps/node_modules/postcss-load-config/src/index.js:55:15)
        at <anonymous>
        at runLoaders (/home/parzibyte/apps/node_modules/webpack/lib/NormalModule.js:195:19)
        at /home/parzibyte/apps/node_modules/loader-runner/lib/LoaderRunner.js:364:11
        at /home/parzibyte/apps/node_modules/loader-runner/lib/LoaderRunner.js:230:18
        at context.callback (/home/parzibyte/apps/node_modules/loader-runner/lib/LoaderRunner.js:111:13)
        at Promise.resolve.then.then.catch (/home/parzibyte/apps/node_modules/postcss-loader/lib/index.js:194:71)
        at <anonymous>
    @ ./node_modules/vuetify/dist/vuetify.min.css
    @ ./src/main.js

    Build failed with errors.

    Gracias a https://github.com/akveo/ngx-admin/issues/604#issuecomment-271974780

*/
module.exports = {};