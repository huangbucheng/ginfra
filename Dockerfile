FROM centos

# Add crontab file in the cron directory
ADD crontab /var/spool/cron/root

# Give execution rights on the cron job
RUN chmod 0644 /var/spool/cron/root

# FOR some Docker-Centos security issue
RUN sed -i '/pam_loginuid.so/ s/^/#/' /etc/pam.d/crond

CMD ["crond", "-n"]
