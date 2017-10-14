FROM scratch
COPY build/tailtopic-linux /tailtopic
ENTRYPOINT [ "/tailtopic" ]
CMD []
