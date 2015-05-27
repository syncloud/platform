from ConfigParser import ConfigParser


class SamConfig:

    def __init__(self, filename, root_dir='', local_dir=''):
        self.local_dir = local_dir
        self.root_dir = root_dir
        self.parser = ConfigParser()
        self.parser.read(filename)
        self.filename = filename

    def save(self):
        with open(self.filename, 'wb') as file:
            self.parser.write(file)

    def section(self, name):
        if not self.parser.has_section(name):
            self.parser.add_section(name)

    def apps_dir(self):
        return self.parser.get('sam', 'apps_dir').format(self.root_dir)

    def set_apps_dir(self, apps_dir):
        self.section('sam')
        self.parser.set('sam', 'apps_dir', apps_dir)
        self.save()

    def status_dir(self):
        return self.parser.get('sam', 'status_dir').format(self.root_dir)

    def set_status_dir(self, status_dir):
        self.section('sam')
        self.parser.set('sam', 'status_dir', status_dir)
        self.save()

    def bin_dir(self):
        return self.parser.get('sam', 'bin_dir').format(self.local_dir)

    def set_bin_dir(self, bin_dir):
        self.section('sam')
        self.parser.set('sam', 'bin_dir', bin_dir)
        self.save()

    def apps_url_template(self):
        if self.parser.has_option('sam', 'apps_url_template'):
            return self.parser.get('sam', 'apps_url_template')
        return None

    def set_apps_url_template(self, apps_url_template):
        self.section('sam')
        self.parser.set('sam', 'apps_url_template', apps_url_template)
        self.save()