FROM golang:1.14.1 as BUILD

ENV Workspace="/covid"
ENV BuildDir="/build"

WORKDIR $Workspace

COPY . .

RUN mkdir $BuildDir \
   && CGO_ENABLED=0 GOOS=linux go build -o "$BuildDir/covid" . \
   && mkdir "$BuildDir/templates" \
   && cp -r ./templates "$BuildDir/templates"

FROM alpine:3.11.3

ENV Workspace="/covid-run"

RUN mkdir $Workspace
WORKDIR $Workspace

COPY  --from=BUILD build/* ./

EXPOSE 5900
CMD ["./covid"]