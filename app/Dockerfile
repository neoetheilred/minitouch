FROM golang as build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

FROM scratch
COPY --from=build /app/main /main
CMD [ "/main" ]