class Applications {
  constructor (settings) {
    this.settings = settings
    this.all = []
  }

  get template () {
    return this.all.filter((app) => {
      return app['is-template']
    })
  }

  get open () {
    return this.all.filter((app) => {
      return app.open
    })
  }

  update (containers) {
    let apps = []
    const enabled = {
      system: this.settings.showSystemApplications,
      template: this.settings.showTemplateApplications
    }
    containers.forEach((container) => {
      if (enabled[container.name] !== undefined && !enabled[container.name]) {
        return
      }
      apps = apps.concat(container.apps || [])
    })
    this.all = apps
  }
}

export default Applications
