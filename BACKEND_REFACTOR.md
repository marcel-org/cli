# Backend modules refactor:

Ok i want to create a new architecture for the backend, because currently it's really messy in the `routes/` folder.

## Problem
In this there is a file for each module. The problem is that each route in those files has:
- validation
- logic,
- routing
in the same scopes, which is crazy....
We should seperate those to have a cleaner architecture.

## New architecture
### Modules
we're going to rename the `routes/` folder to `modules/`  

and now each module inside will have it's own folder. for instance admin folder will be named `admin/`, or focus -> `focus/`

Each module folder will have three files inside: `router.ts`, `service.ts`, `validator.ts`

In the router file we will put the routing and call the validators before calling the service.
The service file will contain the logic and db calls, etc... It must be a class called `AdminService` for admin for example or `CalendarService`, etc...
The validator file will contain all zod validators that are called in the router.

### Modules Loader
Now we will refactor the `webserver.ts`. Instead of importing all the folders by hand we will use a new file called modulesLoader.ts

```ts
import type { Hono } from 'hono';
import { readdirSync, statSync } from 'node:fs';
import { join } from 'node:path';
import { pathToFileURL } from 'node:url';

export async function loadModules(webserver: Hono) {
  const modulesDir = 'src/modules';
  for (const entry of readdirSync(modulesDir)) {
    const moduleDir = join(modulesDir, entry);
    if (!statSync(moduleDir).isDirectory()) continue;

    const moduleRouterFile = join(moduleDir, `router.ts`);
    try {
      const { default: routeHandler } = await import(
        pathToFileURL(moduleRouterFile).href
      );

      webserver.route(`/${entry}`, routeHandler);
      console.log(`[module] ⇢   ${moduleRouterFile}  ->  /${entry}`);
    } catch (err: any) {
      if (err.code === 'ERR_MODULE_NOT_FOUND' || err.code === 'MODULE_NOT_FOUND') {
        console.warn(`[module] (skip) Aucun fichier route pour "${entry}"`);
        console.warn("tried to access but error: ", err.message);
        continue;
      }
      throw err;
    }
  }
}

export default loadModules;
```

Use a logic like this to load all the modules by importing this `modulesLoader.ts` in the `webserver.ts`

With those changes il will enable adding a new module really easily
without having to add again the import in the `webserver.ts` file every time
we add a new module
