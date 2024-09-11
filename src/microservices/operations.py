import requests

def fetch_data_from_api(api_url):
    
    try:
        response = requests.get(api_url)
        response.raise_for_status()  # Raise an HTTPError if the request returned an unsuccessful status code
        data = response.json()  # Assuming the API returns JSON data
        return data.get('datapoints', [])  # Assuming 'datapoints' contains the list of data points
    except requests.exceptions.RequestException as e:
        print(f"Error fetching data from API: {e}")
        return []

def up_or_down(data=None, api_url=None):
    
    # Fetch data from the API if no data is provided
    if data is None and api_url is not None:
        data = fetch_data_from_api(api_url)
    
    # If the data is empty or still None, return 'DOWN'
    if not data:
        return 'DOWN'

    # Count the number of 'UP' (0) and 'DOWN' (1) data points
    up_count = sum(1 for point in data if point[0] == 0)
    total_count = len(data)

    # If total_count is 0 (empty data), return 'DOWN'
    if total_count == 0:
        return 'DOWN'

    # Calculate the percentage of 'UP' points
    up_percentage = (up_count / total_count) * 100

    # If 70% or more of the points are 'UP', return 'UP', otherwise 'DOWN'
    if up_percentage >= 70:
        return 'UP'
    else:
        return 'DOWN'

# Export the function
__all__ = ['up_or_down']
