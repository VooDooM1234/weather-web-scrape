      
      
      function updateWeatherInfo(){
      
      const country = document.getElementById("country");
      const city = document.getElementById("city");
      const temp = document.getElementById("temp");
      const summary = document.getElementById("summary");
      const wind = document.getElementById("wind");
      const pressure = document.getElementById("pressure");
      const precip = document.getElementById("precip");
      const humidity = document.getElementById("humidity");
      const cloud = document.getElementById("cloud");
      const feelslike = document.getElementById("feelslike");
      const windchill = document.getElementById("windchill");
      const heatindex = document.getElementById("heatindex");
      const dewpoint = document.getElementById("dewpoint");
      const visibility = document.getElementById("visibility");
      const uv = document.getElementById("uv");
      const gust = document.getElementById("gust");

      fetch("http://localhost:8080/api/weather")
        .then((response) => {
          if (!response.ok) {
            throw new Error("Network response was not ok " + response.statusText);
          }
          return response.json();
        })
        .then((data) => {
          country.textContent = data.location.country;
          city.textContent = data.location.name;
          temp.textContent = `${data.current.temp_c} °C`;
          summary.textContent = data.current.condition.text;

          wind.textContent = `${data.current.wind_mph} mph (${data.current.wind_kph} kph) ${data.current.wind_dir} (${data.current.wind_degree}°)`;
          pressure.textContent = `${data.current.pressure_mb} mb (${data.current.pressure_in} in)`;
          precip.textContent = `${data.current.precip_mm} mm (${data.current.precip_in} in)`;
          humidity.textContent = `${data.current.humidity}%`;
          cloud.textContent = `${data.current.cloud}%`;
          feelslike.textContent = `${data.current.feelslike_c} °C (${data.current.feelslike_f} °F)`;
          windchill.textContent = `${data.current.windchill_c} °C (${data.current.windchill_f} °F)`;
          heatindex.textContent = `${data.current.heatindex_c} °C (${data.current.heatindex_f} °F)`;
          dewpoint.textContent = `${data.current.dewpoint_c} °C (${data.current.dewpoint_f} °F)`;
          visibility.textContent = `${data.current.vis_km} km (${data.current.vis_miles} miles)`;
          uv.textContent = data.current.uv;
          gust.textContent = `${data.current.gust_mph} mph (${data.current.gust_kph} kph)`;
        })
        .catch((error) => {
          console.error("Fetch error:", error);
        });

    }