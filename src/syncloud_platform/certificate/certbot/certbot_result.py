class CertbotResult:
    def __init__(self, certificate_file, key_file, regenerated):
        self.certificate_file = certificate_file
        self.key_file = key_file
        self.regenerated = regenerated