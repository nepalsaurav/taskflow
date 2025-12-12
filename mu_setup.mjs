#!/usr/bin/env zx

let mu_exit = false

try {
  await $`which mu`
  mu_exit = true
} catch {
  mu_exit = false
}

// install maildir if not only
if (!mu_exit) {
  await $`sudo apt install -y maildir-utils`
  mu_exit = true
}

if (mu_exit) {
  try {
    let index = await $`mu index`
    console.log(index.stdout)
  } catch (err) {
    if (err.stderr.includes("Try (re)creating using `mu init'")) {
      let init = await $`mu init`
      console.log(init.stdout)
      let index = await $`mu index`
      console.log(index.stdout)
    }
  }
}

// create cron tab entry if mu mu_exit

if (mu_exit) {
  let muPath = (await $`which mu`).stdout.trim()
  let existing_cron_entries = ""
  try {
    existing_cron_entries = (await $`crontab -l 2>/dev/null || echo ""`).stdout
  } catch {

  }
  let cronEntry = `*/5 * * * * ${muPath} index`
  if (!existing_cron_entries.includes(cronEntry)) {
    await $`(crontab -l 2>/dev/null; echo "*/5 * * * * ${muPath} index") | crontab -`
    console.log("cron job for mu index successfully")
  } else {
    console.log("cron job already exist")
  }
}
