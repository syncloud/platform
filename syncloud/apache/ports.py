class Ports:
    # These port numbers are in http.conf and https.conf
    def __init__(self, http=80, https=443):
        self.http = http
        self.https = https