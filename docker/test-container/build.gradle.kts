plugins {
    application
}

repositories {
    mavenCentral()
}

dependencies {
    implementation("com.github.org-arl:fjage:2.0.1")
}

application {
    mainClass = "Main"
}
